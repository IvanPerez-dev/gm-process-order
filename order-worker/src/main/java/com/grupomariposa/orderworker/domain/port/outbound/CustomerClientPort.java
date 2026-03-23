package com.grupomariposa.orderworker.domain.port.outbound;

import com.grupomariposa.orderworker.domain.model.Customer;
import reactor.core.publisher.Mono;

public interface CustomerClientPort {
    Mono<Customer> findById(String customerId);
}
