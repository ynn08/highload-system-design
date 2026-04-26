package com.example.orderservice.domain.model;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class Order {
    private UUID orderId;
    private UUID customerId;
    private UUID restaurantId;
    private OrderStatus status;
    private BigDecimal totalAmount;
    private BigDecimal deliveryFee;
    private BigDecimal surgeMultiplier;
    private List<OrderItem> items;
    private String deliveryAddress;
    private LocalDateTime estimatedDelivery;
    private LocalDateTime createdAt;

    public void validate() {
        if (items == null || items.isEmpty()) {
            throw new IllegalStateException("Order must have at least one item");
        }
        if (totalAmount == null || totalAmount.compareTo(BigDecimal.ZERO) <= 0) {
            throw new IllegalStateException("Total amount must be positive");
        }
    }

    public void calculateEstimatedDelivery(int distanceMeters) {
        // Non-trivial logic: 15 mins prep + 2 mins per km
        long minutes = 15 + (distanceMeters / 500);
        this.estimatedDelivery = LocalDateTime.now().plusMinutes(minutes);
    }
}
