package com.grupomariposa.orderworker.infrastructure.mongodb.repository;

import com.grupomariposa.orderworker.infrastructure.mongodb.document.OrderDocument;
import org.springframework.data.mongodb.repository.ReactiveMongoRepository;
import reactor.core.publisher.Mono;

public interface OrderMongoRepository extends ReactiveMongoRepository<OrderDocument, String> {
    //Mono<OrderDocument> findByOrderId(String orderId);
}
