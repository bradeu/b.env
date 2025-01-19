import { contractAPI } from '../../utils/contractAPI';

// REST API endpoint
export default async function handler(req, res) {
    const { address } = req.query;
    
    try {
        const data = await contractAPI.getPublicData(address);
        res.status(200).json({ data });
    } catch (error) {
        res.status(500).json({ error: error.message });
    }
} 