package com.grupomariposa.orderworker.infrastructure.redis.lock;

import com.grupomariposa.orderworker.domain.port.outbound.DistributedLockPort;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.data.redis.core.ReactiveStringRedisTemplate;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.time.Duration;

@Component

public class DistributedLockAdapter implements DistributedLockPort {

    private final ReactiveStringRedisTemplate redisTemplate;

    public DistributedLockAdapter(ReactiveStringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
    }
    private static final String PREFIX = "lock:order:";
    private static final Duration TTL = Duration.ofMinutes(5); // evita lock infinito si el worker muere

    @Override
    public Mono<Boolean> acquire(String orderId) {
        return redisTemplate.opsForValue()
                            .setIfAbsent(PREFIX + orderId, "locked", TTL); // SETNX atomico
    }

    @Override
    public Mono<Void> release(String orderId) {
        return redisTemplate.delete(PREFIX + orderId).then();
    }
}
