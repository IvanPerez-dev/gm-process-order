package com.grupomariposa.orderworker.domain.port.inbound;

import reactor.core.publisher.Mono;

public interface ProcessOrderUseCase {
    Mono<Void> process(OrderMessage message);
}
