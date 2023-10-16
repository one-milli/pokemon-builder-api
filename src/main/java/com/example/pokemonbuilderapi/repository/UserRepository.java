package com.example.pokemonbuilderapi.repository;

import com.example.pokemonbuilderapi.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UserRepository extends JpaRepository<User, Long> {

    User findByUsername(String username);
}
