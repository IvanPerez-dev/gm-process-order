package com.grupomariposa.orderworker.application.usecase;

import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import com.grupomariposa.orderworker.domain.enums.OrderStatus;
import com.grupomariposa.orderworker.domain.model.Customer;
import com.grupomariposa.orderworker.domain.model.Order;
import com.grupomariposa.orderworker.domain.model.OrderItem;
import com.grupomariposa.orderworker.domain.port.inbound.order.OrderMessageCommand;
import com.grupomariposa.orderworker.domain.port.outbound.*;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.test.StepVerifier;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class ProcessOrderUseCaseImplTest {

    @Mock private CustomerClientPort customerPort;
    @Mock private ProductClientPort productPort;
    @Mock private OrderRepository orderRepository;
    @Mock private IdempotencyPort idempotencyPort;
    @Mock private DistributedLockPort distributedLockPort;
    @Mock private RetryStatePort retryStatePort;

    @InjectMocks
    private ProcessOrderUseCaseImpl useCase;

    private static final String ORDER_ID   = "order-001";
    private static final String CUSTOMER_ID = "cust-001";
    private static final String PRODUCT_ID  = "prod-001";

    private OrderMessageCommand command;

    @BeforeEach
    void setUp() {
        command = new OrderMessageCommand(
                ORDER_ID,
                CUSTOMER_ID,
                List.of(new OrderMessageCommand.OrderItemDto(PRODUCT_ID, 2))
        );
    }

    private Order pendingOrder() {
        return Order.builder()
                    .orderId(ORDER_ID)
                    .customer(null)
                    .products(new ArrayList<>(List.of(OrderItem.of(PRODUCT_ID))))
                    .status(OrderStatus.PENDING)
                    .retryCount(0)
                    .errorMsg(null)
                    .createdAt(Instant.now())
                    .processedAt(null)
                    .build();
    }

    private Order processedOrder() {
        return Order.builder()
                    .orderId(ORDER_ID)
                    .customer(null)
                    .products(new ArrayList<>(List.of(OrderItem.of(PRODUCT_ID))))
                    .status(OrderStatus.PROCESSED)
                    .retryCount(0)
                    .errorMsg(null)
                    .createdAt(Instant.now())
                    .processedAt(null)
                    .build();
    }

    private Customer activeCustomer() {
        return Customer.create(CUSTOMER_ID, "Juan Pérez", "juan@test.com", true);
    }

    private OrderItem enrichedItem() {
        OrderItem item = OrderItem.of(PRODUCT_ID);
        item.enrich("Producto Alpha", "Desc Alpha", 99.99);
        return item;
    }

    // -----------------------------------------------------------------------
    // Order not found
    // -----------------------------------------------------------------------

    @Test
    void process_whenOrderNotFound_returnsError() {
        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.empty());

        StepVerifier.create(useCase.process(command))
                    .expectError(OrderProcessingException.class)
                    .verify();
    }

    // -----------------------------------------------------------------------
    // Order not eligible
    // -----------------------------------------------------------------------

    @Test
    void process_whenOrderNotEligible_completesEmpty() {
        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(processedOrder()));

        StepVerifier.create(useCase.process(command))
                    .verifyComplete();

        verifyNoInteractions(distributedLockPort, idempotencyPort, customerPort, productPort);
    }

    // -----------------------------------------------------------------------
    // Lock not acquired
    // -----------------------------------------------------------------------

    @Test
    void process_whenLockNotAcquired_completesEmpty() {
        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(pendingOrder()));
        when(distributedLockPort.acquire(ORDER_ID)).thenReturn(Mono.just(false));

        StepVerifier.create(useCase.process(command))
                    .verifyComplete();

        verifyNoInteractions(idempotencyPort, customerPort, productPort);
    }

    // -----------------------------------------------------------------------
    // Already processed (idempotency)
    // -----------------------------------------------------------------------

    @Test
    void process_whenAlreadyProcessed_completesEmpty() {
        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(pendingOrder()));
        when(distributedLockPort.acquire(ORDER_ID)).thenReturn(Mono.just(true));
        when(distributedLockPort.release(ORDER_ID)).thenReturn(Mono.empty());
        when(idempotencyPort.isAlreadyProcessed(ORDER_ID)).thenReturn(Mono.just(true));

        StepVerifier.create(useCase.process(command))
                    .verifyComplete();

        verifyNoInteractions(customerPort, productPort);
        verify(orderRepository, never()).save(any());
    }

    // -----------------------------------------------------------------------
    // Happy path
    // -----------------------------------------------------------------------

    @Test
    void process_happyPath_enrichesOrderAndMarksAsProcessed() {
        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(pendingOrder()));
        when(distributedLockPort.acquire(ORDER_ID)).thenReturn(Mono.just(true));
        when(distributedLockPort.release(ORDER_ID)).thenReturn(Mono.empty());
        when(idempotencyPort.isAlreadyProcessed(ORDER_ID)).thenReturn(Mono.just(false));
        when(customerPort.findById(CUSTOMER_ID)).thenReturn(Mono.just(activeCustomer()));
        when(productPort.findByIds(List.of(PRODUCT_ID))).thenReturn(Flux.just(enrichedItem()));
        when(orderRepository.save(any())).thenReturn(Mono.empty());
        when(idempotencyPort.markAsProcessed(ORDER_ID)).thenReturn(Mono.empty());

        StepVerifier.create(useCase.process(command))
                    .verifyComplete();

        verify(orderRepository).save(any());
        verify(idempotencyPort).markAsProcessed(ORDER_ID);
    }

    // -----------------------------------------------------------------------
    // Retry below max
    // -----------------------------------------------------------------------

    @Test
    void process_whenEnrichFails_andRetryBelowMax_propagatesError() {
        RuntimeException cause = new RuntimeException("service unavailable");

        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(pendingOrder()));
        when(distributedLockPort.acquire(ORDER_ID)).thenReturn(Mono.just(true));
        when(distributedLockPort.release(ORDER_ID)).thenReturn(Mono.empty());
        when(idempotencyPort.isAlreadyProcessed(ORDER_ID)).thenReturn(Mono.just(false));
        when(customerPort.findById(CUSTOMER_ID)).thenReturn(Mono.error(cause));
        when(retryStatePort.incrementAndGet(ORDER_ID)).thenReturn(Mono.just(1));

        StepVerifier.create(useCase.process(command))
                    .expectError(RuntimeException.class)
                    .verify();

        verify(retryStatePort).incrementAndGet(ORDER_ID);
    }

    // -----------------------------------------------------------------------
    // Retry reaches max → goes to DLQ (completes empty)
    // -----------------------------------------------------------------------

    @Test
    void process_whenEnrichFails_andRetryReachesMax_completesEmpty() {
        RuntimeException cause = new RuntimeException("service unavailable");

        when(orderRepository.findById(ORDER_ID)).thenReturn(Mono.just(pendingOrder()));
        when(distributedLockPort.acquire(ORDER_ID)).thenReturn(Mono.just(true));
        when(distributedLockPort.release(ORDER_ID)).thenReturn(Mono.empty());
        when(idempotencyPort.isAlreadyProcessed(ORDER_ID)).thenReturn(Mono.just(false));
        when(customerPort.findById(CUSTOMER_ID)).thenReturn(Mono.error(cause));
        when(retryStatePort.incrementAndGet(ORDER_ID)).thenReturn(Mono.just(3));

        StepVerifier.create(useCase.process(command))
                    .verifyComplete();

        verify(retryStatePort).incrementAndGet(ORDER_ID);
        verify(orderRepository, never()).save(any());
    }
}
