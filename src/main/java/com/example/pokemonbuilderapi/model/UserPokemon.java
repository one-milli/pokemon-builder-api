package com.example.pokemonbuilderapi.model;

import jakarta.persistence.*;

@Entity
@Table(name="user_pokemons")
public class UserPokemon {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne
    @JoinColumn(name="user_id")
    private User user;

    @ManyToOne
    @JoinColumn(name = "pokemon_id")
    private Pokemon pokemon;

    private String effortValues;
    private String notes;
}
