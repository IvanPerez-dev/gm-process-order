package com.grupomariposa.orderworker.infrastructure.mongodb.mapper;

import com.grupomariposa.orderworker.domain.model.Customer;
import com.grupomariposa.orderworker.domain.model.Order;
import com.grupomariposa.orderworker.domain.model.OrderItem;
import com.grupomariposa.orderworker.infrastructure.mongodb.document.CustomerDocument;
import com.grupomariposa.orderworker.infrastructure.mongodb.document.OrderDocument;
import com.grupomariposa.orderworker.infrastructure.mongodb.document.OrderItemDocument;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface OrderDocumentMapper {
    @Mapping(target = "id", ignore = true)
    @Mapping(target = "status", expression = "java(order.getStatus().name())")
    OrderDocument toDocument(Order order);

    @Mapping(target = "orderId", source = "id")
    @Mapping(target = "status", expression = "java(com.grupomariposa.orderworker.domain.enums.OrderStatus.valueOf" +
            "(document.getStatus()))")
    Order toDomain(OrderDocument document);

    CustomerDocument toCustomerDocument(Customer customer);
    @Mapping(target = "total", expression = "java(item.getTotal())")
    OrderItemDocument toItemDocument(OrderItem item);

    Customer toCustomer(CustomerDocument document);
    OrderItem toOrderItem(OrderItemDocument document);
}
