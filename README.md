# Henchies Backend Go

A backend for a real time multiplayer game written in Go

Designed to be horiontally scalable by relying on Redis for data storage and Photon Networking for player to player synchronization.

Game state is managed on each players device and signed off by backend for other players to avoid potential for modifying local game state. 
This hybrid architecture creates an easy to scale system with minimal runtime costs.
