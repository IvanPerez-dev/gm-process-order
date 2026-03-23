package com.grupomariposa.orderworker.infrastructure.redis.retry;

import com.grupomariposa.orderworker.domain.port.outbound.RetryStatePort;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.ReactiveRedisTemplate;
import org.springframework.data.redis.core.ReactiveStringRedisTemplate;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.time.Duration;

@Component
public class RetryStateAdapter  implements RetryStatePort {

    private final ReactiveStringRedisTemplate redisTemplate;

    public RetryStateAdapter(ReactiveStringRedisTemplate redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    private static final String PREFIX = "retry:order:";
    private static final Duration TTL = Duration.ofHours(24);

    @Override
    public Mono<Integer> incrementAndGet(String orderId) {
        String key = PREFIX + orderId;
        return redisTemplate.opsForValue()
                            .increment(key)
                            .flatMap(count ->
                                             redisTemplate.expire(key, TTL)
                                                          .thenReturn(count.intValue())
                            );
    }

    @Override
    public Mono<Void> clear(String orderId) {
        return redisTemplate.delete(PREFIX + orderId).then();
    }
}
