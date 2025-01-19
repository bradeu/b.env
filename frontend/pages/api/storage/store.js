import { API } from '../../../utils/contractAPI';

export default async function handler(req, res) {
    if (req.method !== 'POST') {
        return res.status(405).json({ error: 'Method not allowed' });
    }

    const { address, encryptedKey, proof } = req.body;
    
    try {
        const result = await API.storeApiKey(address, encryptedKey, proof);
        if (result.success) {
            res.status(200).json(result);
        } else {
            res.status(400).json(result);
        }
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
} 