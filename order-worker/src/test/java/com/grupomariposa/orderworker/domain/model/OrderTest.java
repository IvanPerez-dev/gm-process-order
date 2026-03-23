package com.grupomariposa.orderworker.domain.model;

import com.grupomariposa.orderworker.domain.ecexption.CustomerNotActiveException;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import com.grupomariposa.orderworker.domain.enums.OrderStatus;
import org.junit.jupiter.api.Test;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

import static org.assertj.core.api.Assertions.*;

class OrderTest {

    private Order buildOrder(OrderStatus status, List<OrderItem> items) {
        return Order.builder()
                    .orderId("order-001")
                    .customer(null)
                    .products(new ArrayList<>(items))
                    .status(status)
                    .retryCount(0)
                    .errorMsg(null)
                    .createdAt(Instant.now())
                    .processedAt(null)
                    .build();
    }

    private Customer activeCustomer() {
        return Customer.create("cust-1", "Juan Pérez", "juan@test.com", true);
    }

    private Customer inactiveCustomer() {
        return Customer.create("cust-2", "María García", "maria@test.com", false);
    }

    private OrderItem enrichedItem(String productId) {
        OrderItem item = OrderItem.of(productId);
        item.enrich("Producto A", "Descripción A", 99.99);
        return item;
    }


    @Test
    void isEligibleForProcessing_whenPending_returnsTrue() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));
        assertThat(order.isEligibleForProcessing()).isTrue();
    }

    @Test
    void isEligibleForProcessing_whenProcessed_returnsFalse() {
        Order order = buildOrder(OrderStatus.PROCESSED, List.of(OrderItem.of("p-1")));
        assertThat(order.isEligibleForProcessing()).isFalse();
    }

    @Test
    void isEligibleForProcessing_whenFailed_returnsFalse() {
        Order order = buildOrder(OrderStatus.FAILED, List.of(OrderItem.of("p-1")));
        assertThat(order.isEligibleForProcessing()).isFalse();
    }

    @Test
    void isEligibleForProcessing_whenProcessing_returnsFalse() {
        Order order = buildOrder(OrderStatus.PROCESSING, List.of(OrderItem.of("p-1")));
        assertThat(order.isEligibleForProcessing()).isFalse();
    }

    // --- markAsProcessing ---

    @Test
    void markAsProcessing_whenPending_changesStatusToProcessing() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));
        order.markAsProcessing();
        assertThat(order.getStatus()).isEqualTo(OrderStatus.PROCESSING);
    }

    @Test
    void markAsProcessing_whenNotPending_throwsOrderProcessingException() {
        Order order = buildOrder(OrderStatus.PROCESSED, List.of(OrderItem.of("p-1")));
        assertThatThrownBy(order::markAsProcessing)
                .isInstanceOf(OrderProcessingException.class)
                .hasMessageContaining("Invalid transition");
    }

    // --- enrich ---

    @Test
    void enrich_withActiveCustomerAndValidItems_setsProcessedStatusAndUpdatesProducts() {
        OrderItem raw = OrderItem.of("p-1");
        Order order = buildOrder(OrderStatus.PENDING, List.of(raw));
        OrderItem enriched = enrichedItem("p-1");

        order.enrich(activeCustomer(), List.of(enriched));

        assertThat(order.getStatus()).isEqualTo(OrderStatus.PROCESSED);
        assertThat(order.getCustomer()).isNotNull();
        assertThat(order.getProcessedAt()).isNotNull();
        assertThat(order.getProducts()).hasSize(1);
        assertThat(order.getProducts().get(0).isEnriched()).isTrue();
    }

    @Test
    void enrich_withInactiveCustomer_throwsCustomerNotActiveException() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));
        OrderItem enriched = enrichedItem("p-1");

        assertThatThrownBy(() -> order.enrich(inactiveCustomer(), List.of(enriched)))
                .isInstanceOf(CustomerNotActiveException.class)
                .hasMessageContaining("is not active");
    }

    @Test
    void enrich_withNullItems_throwsOrderProcessingException() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));

        assertThatThrownBy(() -> order.enrich(activeCustomer(), null))
                .isInstanceOf(OrderProcessingException.class)
                .hasMessageContaining("empty items");
    }

    @Test
    void enrich_withEmptyItems_throwsOrderProcessingException() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));

        assertThatThrownBy(() -> order.enrich(activeCustomer(), List.of()))
                .isInstanceOf(OrderProcessingException.class)
                .hasMessageContaining("empty items");
    }

    @Test
    void enrich_withMismatchedItemCount_throwsOrderProcessingException() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1"), OrderItem.of("p-2")));
        List<OrderItem> enrichedItems = List.of(enrichedItem("p-1"));

        assertThatThrownBy(() -> order.enrich(activeCustomer(), enrichedItems))
                .isInstanceOf(OrderProcessingException.class)
                .hasMessageContaining("mismatch");
    }

    @Test
    void enrich_withNotEnrichedItems_throwsOrderProcessingException() {
        Order order = buildOrder(OrderStatus.PENDING, List.of(OrderItem.of("p-1")));
        OrderItem notEnriched = OrderItem.of("p-1");

        assertThatThrownBy(() -> order.enrich(activeCustomer(), List.of(notEnriched)))
                .isInstanceOf(OrderProcessingException.class)
                .hasMessageContaining("Not all items were enriched");
    }
}
