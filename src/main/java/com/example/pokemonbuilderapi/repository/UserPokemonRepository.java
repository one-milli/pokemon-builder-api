package com.example.pokemonbuilderapi.repository;

import com.example.pokemonbuilderapi.model.UserPokemon;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface UserPokemonRepository extends JpaRepository<UserPokemon, Long> {

    List<UserPokemon> findByUserId(Long userId);
}
