package com.grupomariposa.orderworker.domain.model;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class OrderItemTest {

    @Test
    void of_createsUnenrichedItem() {
        OrderItem item = OrderItem.of("prod-1");

        assertThat(item.getProductId()).isEqualTo("prod-1");
        assertThat(item.getName()).isNull();
        assertThat(item.getPrice()).isNull();
        assertThat(item.isEnriched()).isFalse();
    }

    @Test
    void enrich_setsNameDescriptionAndPrice() {
        OrderItem item = OrderItem.of("prod-1");

        item.enrich("Producto Alpha", "Descripción Alpha", 149.99);

        assertThat(item.getName()).isEqualTo("Producto Alpha");
        assertThat(item.getDescription()).isEqualTo("Descripción Alpha");
        assertThat(item.getPrice()).isEqualTo(149.99);
    }

    @Test
    void isEnriched_whenNameAndPricePresent_returnsTrue() {
        OrderItem item = OrderItem.of("prod-1");
        item.enrich("Nombre", "Desc", 50.0);

        assertThat(item.isEnriched()).isTrue();
    }

    @Test
    void isEnriched_whenNameIsNull_returnsFalse() {
        OrderItem item = OrderItem.of("prod-1");
        item.enrich(null, "Desc", 50.0);

        assertThat(item.isEnriched()).isFalse();
    }

    @Test
    void isEnriched_whenPriceIsNull_returnsFalse() {
        OrderItem item = OrderItem.of("prod-1");
        item.enrich("Nombre", "Desc", null);

        assertThat(item.isEnriched()).isFalse();
    }

    @Test
    void isEnriched_whenBothNullAfterOf_returnsFalse() {
        OrderItem item = OrderItem.of("prod-99");
        assertThat(item.isEnriched()).isFalse();
    }
}
