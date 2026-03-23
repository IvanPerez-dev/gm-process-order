package com.grupomariposa.orderworker.domain.port.outbound;

import reactor.core.publisher.Mono;

public interface IdempotencyPort {
    Mono<Boolean> isAlreadyProcessed(String orderId);
    Mono<Void> markAsProcessed(String orderId);
}
