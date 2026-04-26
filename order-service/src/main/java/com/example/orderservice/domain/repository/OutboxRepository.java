package com.example.orderservice.domain.repository;

import com.example.orderservice.domain.model.OutboxEvent;
import java.util.List;

public interface OutboxRepository {
    void save(OutboxEvent event);
    List<OutboxEvent> findPending();
    void markAsProcessed(OutboxEvent event);
}
