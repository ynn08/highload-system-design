package com.example.orderservice.infrastructure.persistence;

import jakarta.persistence.*;
import lombok.Data;
import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "orders")
@Data
public class OrderEntity {
    @Id
    private UUID orderId;
    private UUID customerId;
    private UUID restaurantId;
    private String status;
    private BigDecimal totalAmount;
    private BigDecimal deliveryFee;
    private BigDecimal surgeMultiplier;
    private String deliveryAddress;
    private LocalDateTime estimatedDelivery;
    private LocalDateTime createdAt;
}
