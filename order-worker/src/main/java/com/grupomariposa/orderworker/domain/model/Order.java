package com.grupomariposa.orderworker.domain.model;

import com.grupomariposa.orderworker.domain.enums.OrderStatus;
import com.grupomariposa.orderworker.domain.ecexption.CustomerNotActiveException;
import com.grupomariposa.orderworker.domain.ecexption.OrderProcessingException;
import lombok.*;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@Getter
@Builder
@AllArgsConstructor(access = AccessLevel.PRIVATE)
public class Order {
    private final String orderId;
    private Customer customer;
    private final List<OrderItem> products;
    private OrderStatus status;
    private int retryCount;
    private String errorMsg;
    private final Instant createdAt;
    private Instant processedAt;


    public void markAsProcessing() {
        validateTransition(OrderStatus.PENDING);
        this.status = OrderStatus.PROCESSING;
    }

    public void enrich(Customer customer, List<OrderItem> enrichedItems) {
        if (!customer.getIsActive())
            throw new CustomerNotActiveException(orderId, customer.getId());

        if (enrichedItems == null || enrichedItems.isEmpty())
            throw new OrderProcessingException(orderId, "Cannot enrich order with empty items");

        if (enrichedItems.size() != this.products.size())
            throw new OrderProcessingException(orderId,
                                               "Enriched items count mismatch: expected " + this.products.size()
                                                       + " got " + enrichedItems.size());

        boolean allEnriched = enrichedItems.stream().allMatch(OrderItem::isEnriched);
        if (!allEnriched)
            throw new OrderProcessingException(orderId, "Not all items were enriched");

        Map<String, Integer> quantityMap = this.products.stream()
                                                        .collect(Collectors.toMap(
                                                                OrderItem::getProductId,
                                                                OrderItem::getQuantity
                                                        ));


        enrichedItems.forEach(item -> {
            Integer quantity = quantityMap.get(item.getProductId());
            if (quantity == null) {
                throw new OrderProcessingException(orderId,
                                                   "Product not found in original order: " + item.getProductId());
            }
            item.setQuantity(quantity);
        });

        this.products.clear();
        this.products.addAll(enrichedItems);
        this.customer = customer;
        this.status = OrderStatus.PROCESSED;
        this.processedAt = Instant.now();
    }

    public boolean isEligibleForProcessing() {
        return this.status == OrderStatus.PENDING;
    }

    private void validateTransition(OrderStatus requiredStatus) {
        if (this.status != requiredStatus)
            throw new OrderProcessingException(orderId,
                                               "Invalid transition: expected " + requiredStatus + " but was " + this.status);
    }
}
