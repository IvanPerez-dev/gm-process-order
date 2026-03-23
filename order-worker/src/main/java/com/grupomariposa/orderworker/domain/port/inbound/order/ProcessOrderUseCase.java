package com.grupomariposa.orderworker.domain.port.inbound.order;

import reactor.core.publisher.Mono;

public interface ProcessOrderUseCase {
    Mono<Void> process(OrderMessageCommand message);
}
