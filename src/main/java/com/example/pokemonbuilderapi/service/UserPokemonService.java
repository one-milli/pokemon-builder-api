package com.example.pokemonbuilderapi.service;

import com.example.pokemonbuilderapi.model.UserPokemon;
import com.example.pokemonbuilderapi.repository.UserPokemonRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class UserPokemonService {

    @Autowired
    private UserPokemonRepository userPokemonRepository;

    public List<UserPokemon> getUserPokemonsByUserId(Long userId) {
        return userPokemonRepository.findByUserId(userId);
    }

    public UserPokemon getUserPokemonById(Long id) {
        return userPokemonRepository.findById(id).orElse(null);
    }

    public UserPokemon saveUserPokemon(UserPokemon userPokemon) {
        return userPokemonRepository.save(userPokemon);
    }

    public void deleteUserPokemon(Long id) {
        userPokemonRepository.deleteById(id);
    }
}
