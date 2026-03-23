package com.grupomariposa.orderworker.domain.port.inbound.order;

import java.util.List;

public record OrderMessageCommand(String orderId, String customerId, List<OrderItemDto> items) {
    public record OrderItemDto(String productId, int Quantity){}
}


