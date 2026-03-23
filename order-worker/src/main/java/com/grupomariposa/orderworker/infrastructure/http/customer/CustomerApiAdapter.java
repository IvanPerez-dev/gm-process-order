package com.grupomariposa.orderworker.infrastructure.http.customer;

import com.grupomariposa.orderworker.domain.model.Customer;
import com.grupomariposa.orderworker.domain.port.outbound.CustomerClientPort;
import com.grupomariposa.orderworker.domain.ecexption.CustomerNotFoundException;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpStatusCode;
import org.springframework.stereotype.Component;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;
import reactor.util.retry.Retry;

import java.time.Duration;

@Slf4j
@Component
public class CustomerApiAdapter implements CustomerClientPort {

    private final WebClient webClient;

    public CustomerApiAdapter(@Qualifier("customerWebClient") WebClient webClient) {
        this.webClient = webClient;
    }

    @Override
    public Mono<Customer> findById(String customerId) {
        return webClient.get()
                        .uri("/api/v1/customers/{id}", customerId)
                        .retrieve()
                        .onStatus(HttpStatusCode::is4xxClientError, response ->
                                response.statusCode().equals(HttpStatus.NOT_FOUND)
                                        ? Mono.error(new CustomerNotFoundException(customerId))
                                        : Mono.error(new OrderProcessingException(customerId, "Client error"))
                        )
                        .onStatus(HttpStatusCode::is5xxServerError, response ->
                                Mono.error(new OrderProcessingException(customerId, "Go API server error"))
                        )
                        .bodyToMono(CustomerApiResponse.class)
                        .map(this::toDomain)
                        .retryWhen(Retry.backoff(3, Duration.ofSeconds(2))
                                        .maxBackoff(Duration.ofSeconds(10))
                                        .filter(ex -> ex instanceof OrderProcessingException)
                                        .doBeforeRetry(signal ->
                                                               log.warn("Retrying customer API call, attempt: {}", signal.totalRetries())
                                        )
                        );
    }

    private Customer toDomain(CustomerApiResponse response) {
        return  Customer.create(
                response.id(),
                response.name(),
                response.email(),
                response.isActive()
        );
    }
}
