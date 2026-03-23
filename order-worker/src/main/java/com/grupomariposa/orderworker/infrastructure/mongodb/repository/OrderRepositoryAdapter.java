package com.grupomariposa.orderworker.infrastructure.mongodb.repository;

import com.grupomariposa.orderworker.domain.model.Order;
import com.grupomariposa.orderworker.domain.port.outbound.OrderRepository;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import com.grupomariposa.orderworker.infrastructure.mongodb.mapper.OrderDocumentMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderRepositoryAdapter implements OrderRepository {
    private final OrderMongoRepository mongoRepository;
    private final OrderDocumentMapper orderDocumentMapper;

    @Override
    public Mono<Order> findById(String orderId) {
        return mongoRepository.findById(orderId)
                              .map(orderDocumentMapper::toDomain)
                              .switchIfEmpty(Mono.error(
                                      new OrderProcessingException(orderId, "Order not found in MongoDB")
                              ));
    }


    @Override
    public Mono<Void> save(Order order) {
        return mongoRepository.findById(order.getOrderId())
                              .map(existingDoc -> {

                                  var updated = orderDocumentMapper.toDocument(order);
                                  updated.setId(existingDoc.getId());
                                  return updated;
                              })
                              .switchIfEmpty(Mono.fromSupplier(() -> orderDocumentMapper.toDocument(order)))
                              .flatMap(mongoRepository::save)
                              .doOnSuccess(doc ->
                                                   log.info("Order {} saved with status {}", doc.getId(), doc.getStatus())
                              )
                              .then();
    }
}
