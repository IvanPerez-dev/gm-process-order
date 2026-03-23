package com.grupomariposa.orderworker.domain.port.outbound;

import com.grupomariposa.orderworker.domain.model.OrderItem;
import reactor.core.publisher.Flux;

import java.util.List;

public interface ProductClientPort {
    Flux<OrderItem> findByIds(List<String> productIds);
}
