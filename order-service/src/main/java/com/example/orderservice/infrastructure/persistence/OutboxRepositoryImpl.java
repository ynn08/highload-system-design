package com.example.orderservice.infrastructure.persistence;

import com.example.orderservice.domain.model.OutboxEvent;
import com.example.orderservice.domain.repository.OutboxRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Component;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

@Repository
interface JpaOutboxRepository extends JpaRepository<OutboxEventEntity, UUID> {
    List<OutboxEventEntity> findByStatus(String status);
}

@Component
@RequiredArgsConstructor
public class OutboxRepositoryImpl implements OutboxRepository {
    private final JpaOutboxRepository jpaOutboxRepository;

    @Override
    public void save(OutboxEvent event) {
        OutboxEventEntity entity = new OutboxEventEntity();
        entity.setId(event.getId());
        entity.setAggregateId(event.getAggregateId());
        entity.setAggregateType(event.getAggregateType());
        entity.setEventType(event.getEventType());
        entity.setPayload(event.getPayload());
        entity.setStatus(event.getStatus());
        entity.setCreatedAt(event.getCreatedAt());
        jpaOutboxRepository.save(entity);
    }

    @Override
    public List<OutboxEvent> findPending() {
        return jpaOutboxRepository.findByStatus("PENDING").stream()
                .map(entity -> OutboxEvent.builder()
                        .id(entity.getId())
                        .aggregateId(entity.getAggregateId())
                        .aggregateType(entity.getAggregateType())
                        .eventType(entity.getEventType())
                        .payload(entity.getPayload())
                        .status(entity.getStatus())
                        .createdAt(entity.getCreatedAt())
                        .build())
                .collect(Collectors.toList());
    }

    @Override
    public void markAsProcessed(OutboxEvent event) {
        jpaOutboxRepository.findById(event.getId()).ifPresent(entity -> {
            entity.setStatus("PROCESSED");
            jpaOutboxRepository.save(entity);
        });
    }
}
