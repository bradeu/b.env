// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract SecretStorage is Ownable, ReentrancyGuard {
    struct Secret {
        string encryptedContent;
        address owner;
        bool exists;
    }

    mapping(bytes32 => Secret) private secrets;
    
    event SecretStored(bytes32 indexed secretId, address indexed owner);
    event SecretRetrieved(bytes32 indexed secretId, address indexed retriever);

    constructor() Ownable(msg.sender) {}

    function storeSecret(bytes32 secretId, string memory encryptedContent) 
        external 
        nonReentrant 
    {
        require(!secrets[secretId].exists, "Secret ID already exists");
        
        secrets[secretId] = Secret({
            encryptedContent: encryptedContent,
            owner: msg.sender,
            exists: true
        });

        emit SecretStored(secretId, msg.sender);
    }

    function getSecret(bytes32 secretId) 
        external 
        view 
        returns (string memory) 
    {
        Secret memory secret = secrets[secretId];
        require(secret.exists, "Secret does not exist");
        require(
            secret.owner == msg.sender,
            "Only the owner can retrieve the secret"
        );
        
        return secret.encryptedContent;
    }

    function doesSecretExist(bytes32 secretId) 
        external 
        view 
        returns (bool) 
    {
        return secrets[secretId].exists;
    }
} 