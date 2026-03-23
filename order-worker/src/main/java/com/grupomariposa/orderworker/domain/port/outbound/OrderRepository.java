package com.grupomariposa.orderworker.domain.port.outbound;

import com.grupomariposa.orderworker.domain.enums.OrderStatus;
import com.grupomariposa.orderworker.domain.model.Order;
import reactor.core.publisher.Mono;

public interface OrderRepository {

    Mono<Order> findById(String orderId);
    Mono<Void> save(Order order);
}
