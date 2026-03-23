package com.grupomariposa.orderworker.domain.port.outbound;

import reactor.core.publisher.Mono;

public interface DistributedLockPort {
    Mono<Boolean> acquire(String orderId);
    Mono<Void> release(String orderId);
}
