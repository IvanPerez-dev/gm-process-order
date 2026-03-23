package com.grupomariposa.orderworker.infrastructure.kafka.consumer;

import com.grupomariposa.orderworker.domain.port.inbound.order.OrderMessageCommand;
import com.grupomariposa.orderworker.domain.port.inbound.order.ProcessOrderUseCase;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;
import org.springframework.boot.ApplicationArguments;
import org.springframework.boot.ApplicationRunner;
import org.springframework.kafka.core.reactive.ReactiveKafkaConsumerTemplate;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.util.retry.Retry;

import java.time.Duration;

@Slf4j
@Component
@RequiredArgsConstructor
public class OrderKafkaConsumer   implements ApplicationRunner {
    private final ReactiveKafkaConsumerTemplate<String, OrderMessageCommand> kafkaTemplate;
    private final ProcessOrderUseCase processOrderUseCase;

    @Override
    public void run(ApplicationArguments args) {
        consume()
                .doOnError(ex -> log.error("Fatal error in kafka consumer", ex))
                .retryWhen(Retry.backoff(Long.MAX_VALUE, Duration.ofSeconds(5))
                                .maxBackoff(Duration.ofSeconds(30))
                                // reintenta indefinidamente si el consumer cae
                                .doBeforeRetry(signal ->
                                                       log.warn("Retrying kafka consumer after failure, attempt: {}",
                                                                signal.totalRetries())
                                )
                )
                .subscribe();
    }

    private Flux<Void> consume() {
        return kafkaTemplate.receiveAutoAck()
                            .flatMap(record -> {
                                MDC.put("orderId", record.key());
                                MDC.put("topic", record.topic());
                                MDC.put("partition", String.valueOf(record.partition()));
                                MDC.put("offset", String.valueOf(record.offset()));

                                log.info("Message received — orderId: {}", record.key());

                                return processOrderUseCase.process(record.value())
                                                          .doOnSuccess(v ->
                                                                               log.info("Message processed — orderId: {}", record.key())
                                                          )
                                                          .doOnError(ex ->
                                                                             log.error("Message failed — orderId: {}, error: {}",
                                                                                       record.key(), ex.getMessage())
                                                          )
                                                          .onErrorResume(ex -> Mono.empty()) // no detiene el Flux
                                                          .doFinally(signal -> MDC.clear());
                            },10);
    }
}
