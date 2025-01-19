// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

contract SecretStorage {
    struct ApiLots {
        address lotAddress;
        string encryptedApiKey;
    }

    mapping(address => ApiLots) public users;

    event EncryptedApiKeyStored(address indexed user, string encryptedApiKey);

    function storeEncryptedApiKey(string memory _encryptedApiKey) public {
        require(
            bytes(users[msg.sender].encryptedApiKey).length == 0,
            "API key already stored"
        );

        users[msg.sender] = ApiLots({
            lotAddress: msg.sender,
            encryptedApiKey: _encryptedApiKey
        });

        emit EncryptedApiKeyStored(msg.sender, _encryptedApiKey);
    }

    function getEncryptedApiKey() public view returns (string memory) {
        return users[msg.sender].encryptedApiKey;
    }
}
