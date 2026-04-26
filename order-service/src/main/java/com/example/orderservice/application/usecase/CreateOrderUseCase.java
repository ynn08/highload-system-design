package com.example.orderservice.application.usecase;

import com.example.orderservice.domain.model.Order;
import com.example.orderservice.domain.model.OrderStatus;
import com.example.orderservice.domain.model.OutboxEvent;
import com.example.orderservice.domain.repository.CartRepository;
import com.example.orderservice.domain.repository.OrderRepository;
import com.example.orderservice.domain.repository.OutboxRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.time.LocalDateTime;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class CreateOrderUseCase {
    private final OrderRepository orderRepository;
    private final OutboxRepository outboxRepository;
    private final CartRepository cartRepository;
    private final ObjectMapper objectMapper;

    @Transactional
    @SneakyThrows
    public Order execute(Order order) {
        order.setOrderId(UUID.randomUUID());
        order.setStatus(OrderStatus.PENDING);
        order.setCreatedAt(LocalDateTime.now());
        
        // 1. Dynamic Pricing Logic
        calculateSurgeAndFees(order);
        
        // 2. ETA calculation
        order.calculateEstimatedDelivery(5000); // simulated 5km
        
        order.validate();
        
        // 3. Persist Order
        Order savedOrder = orderRepository.save(order);
        
        // 4. Transactional Outbox: Save event instead of publishing to Kafka
        OutboxEvent outboxEvent = OutboxEvent.builder()
                .id(UUID.randomUUID())
                .aggregateId(savedOrder.getOrderId())
                .aggregateType("ORDER")
                .eventType("ORDER_CREATED")
                .payload(objectMapper.writeValueAsString(savedOrder))
                .status("PENDING")
                .createdAt(LocalDateTime.now())
                .build();
        outboxRepository.save(outboxEvent);
        
        // 5. Clear cart
        cartRepository.deleteByCustomerId(order.getCustomerId());
        
        return savedOrder;
    }

    private void calculateSurgeAndFees(Order order) {
        // Simple Surge logic: multiplier based on hour
        int hour = LocalDateTime.now().getHour();
        BigDecimal multiplier = BigDecimal.ONE;
        
        // Peak hours: 12-14 and 18-21
        if ((hour >= 12 && hour <= 14) || (hour >= 18 && hour <= 21)) {
            multiplier = new BigDecimal("1.5");
        }
        
        BigDecimal baseFee = new BigDecimal("5.00"); // Base 5 units
        BigDecimal distanceKm = new BigDecimal("5"); // Simulated distance
        BigDecimal perKmRate = new BigDecimal("2.00");
        
        BigDecimal fee = baseFee.add(distanceKm.multiply(perKmRate)).multiply(multiplier);
        
        order.setSurgeMultiplier(multiplier);
        order.setDeliveryFee(fee.setScale(2, RoundingMode.HALF_UP));
    }
}
