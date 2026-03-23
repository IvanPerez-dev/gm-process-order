package com.grupomariposa.orderworker.application.usecase;

import com.grupomariposa.orderworker.domain.model.Order;
import com.grupomariposa.orderworker.domain.port.inbound.order.OrderMessageCommand;
import com.grupomariposa.orderworker.domain.port.inbound.order.ProcessOrderUseCase;
import com.grupomariposa.orderworker.domain.port.outbound.*;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

@Slf4j
@Component
@RequiredArgsConstructor
public class ProcessOrderUseCaseImpl implements ProcessOrderUseCase {


    private final CustomerClientPort customerPort;
    private final ProductClientPort productPort;
    private final OrderRepository orderRepository;
    private final IdempotencyPort idempotencyPort;
    private final DistributedLockPort distributedLockPort;
    private final RetryStatePort retryStatePort;

    private static final int MAX_RETRIES = 3;
    @Override
    public Mono<Void> process(OrderMessageCommand message) {
        return orderRepository.findById(message.orderId())
                              .switchIfEmpty(Mono.error( new OrderProcessingException(message.orderId(), "Order not found")))
                              .flatMap( order -> {
                                  if(!order.isEligibleForProcessing()){
                                      log.warn("Order {} is not eligible, current status: {}",
                                               order.getOrderId(), order.getStatus().name());
                                      return Mono.empty();
                                  }
                                  return acquireLockAndProcess(order, message);
                              });
    }

    private Mono<Void> acquireLockAndProcess(Order order, OrderMessageCommand command){
        return distributedLockPort.acquire(command.orderId())
              .flatMap(locked -> {
                  if(!locked){
                      log.warn("Order {} already being processed by another instance",
                               command.orderId());
                      return Mono.empty();
                  }
                return checkIdempotencyAndProcess(order, command)
                        .doFinally(signal -> {
                            distributedLockPort.release(command.orderId())
                                               .doOnError(ex ->
                                                                  log.error("Failed to release lock for order {}", command.orderId(), ex)
                                               )
                                               .subscribe();
            });
         });
    }

    private Mono<Void> checkIdempotencyAndProcess(Order order, OrderMessageCommand command){
        return idempotencyPort.isAlreadyProcessed(command.orderId())
                .flatMap(alreadyProcessed -> {
                    if(alreadyProcessed){
                        log.info("Order {} already processed, skipping", command.orderId());
                        return Mono.empty();
                    }
                    return enrich(command, order);
        });
    }
    private Mono<Void> enrich(OrderMessageCommand command, Order order){
        var productIds = command.items().stream().map(OrderMessageCommand.OrderItemDto::productId).toList();
        return Mono.zip(
                customerPort.findById(command.customerId()),
                productPort.findByIds(productIds).collectList()
        ).flatMap(tuple ->{
            var customer = tuple.getT1();
            var products = tuple.getT2();

            order.enrich(customer, products);

            return orderRepository.save(order);
        }).then(
                idempotencyPort.markAsProcessed(command.orderId())
        ).doOnSuccess(v ->
                log.info("Order {} processed successfully", command.orderId())
        ).onErrorResume(ex -> retryStatePort.incrementAndGet(command.orderId())
          .flatMap(retryCount -> {
              if (retryCount >= MAX_RETRIES) {
                  log.error("Order {} exceeded max retries, sending to DLQ", command.orderId());

                  return Mono.empty();
              }
              log.warn("Order {} failed, retry {}/{}", command.orderId(), retryCount, MAX_RETRIES);
              return Mono.error(ex);
          })
          );
    }
}
