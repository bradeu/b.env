import { storeApiKey, getApiKey } from '../utils/contractInteraction';

export default function ApiKeyManager() {
    async function handleStore() {
        try {
            await storeApiKey(targetAddress, encryptedKey, proof);
            console.log("Stored successfully!");
        } catch (error) {
            console.error("Error storing:", error);
        }
    }

    async function handleRetrieve() {
        try {
            const key = await getApiKey(targetAddress, proof);
            console.log("Retrieved:", key);
        } catch (error) {
            console.error("Error retrieving:", error);
        }
    }

    return (
        <div>
            <button onClick={handleStore}>Store API Key</button>
            <button onClick={handleRetrieve}>Get API Key</button>
        </div>
    );
} 