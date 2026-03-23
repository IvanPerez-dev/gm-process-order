package com.grupomariposa.orderworker.infrastructure.http.product;

import com.grupomariposa.orderworker.domain.model.OrderItem;
import com.grupomariposa.orderworker.domain.port.outbound.ProductClientPort;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatusCode;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.util.retry.Retry;

import java.time.Duration;
import java.util.List;

@Slf4j
@Component
public class ProductApiAdapter implements ProductClientPort {
    private final WebClient webClient;

    public ProductApiAdapter(@Qualifier("productWebClient") WebClient webClient) {
        this.webClient = webClient;
    }

    @Override
    public Flux<OrderItem> findByIds(List<String> productIds) {
        return webClient.get()
                        .uri(uriBuilder -> uriBuilder
                                .path("/api/v1/products")
                                .queryParam("ids", String.join(",", productIds))
                                .build()
                        )
                        .retrieve()
                        .onStatus(HttpStatusCode::is4xxClientError, response ->
                                Mono.error(new OrderProcessingException("products", "Product not found in catalog"))
                        )
                        .onStatus(HttpStatusCode::is5xxServerError, response ->
                                Mono.error(new OrderProcessingException("products", "Go API server error"))
                        )
                        .bodyToFlux(ProductApiResponse.class)
                        .map(this::toDomain)
                        .retryWhen(Retry.backoff(3, Duration.ofSeconds(2))
                                        .maxBackoff(Duration.ofSeconds(10))
                                        .filter(ex -> ex instanceof OrderProcessingException)
                                        .doBeforeRetry(signal ->
                                                               log.warn("Retrying product API call, attempt: {}", signal.totalRetries())
                                        )
                        );
    }

    private OrderItem toDomain(ProductApiResponse response) {
        var item = OrderItem.of(response.id());
        item.enrich(response.name(), response.description(), response.price());
        return item;
    }
}
