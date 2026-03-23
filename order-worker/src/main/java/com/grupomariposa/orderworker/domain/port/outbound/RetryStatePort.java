package com.grupomariposa.orderworker.domain.port.outbound;

import reactor.core.publisher.Mono;

public interface RetryStatePort {
    Mono<Integer> incrementAndGet(String orderId);
    Mono<Void> clear(String orderId);
}
