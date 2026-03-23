package com.grupomariposa.orderworker.infrastructure.redis.idempotency;

import com.grupomariposa.orderworker.domain.port.outbound.IdempotencyPort;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.data.redis.core.ReactiveStringRedisTemplate;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.time.Duration;

@Component

public class IdempotencyAdapter implements IdempotencyPort {
    private final ReactiveStringRedisTemplate redisTemplate;

    public IdempotencyAdapter(ReactiveStringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    private static final String PREFIX = "idempotency:order:";
    private static final Duration TTL = Duration.ofHours(24);

    @Override
    public Mono<Boolean> isAlreadyProcessed(String orderId) {
        return redisTemplate.hasKey(PREFIX + orderId);
    }

    @Override
    public Mono<Void> markAsProcessed(String orderId) {
        return redisTemplate.opsForValue()
                            .set(PREFIX + orderId, "processed", TTL)
                            .then();
    }
}
