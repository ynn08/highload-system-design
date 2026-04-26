package com.example.orderservice.domain.repository;

import com.example.orderservice.domain.model.Order;

public interface MessagePublisher {
    void publishOrderCreated(Order order);
}
