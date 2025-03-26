interface GameSessionResponse {
    session_id: string;
  }
  
export async function createGameSession(playerNames: string[]): Promise<string> {
    const requestHeaders: HeadersInit = new Headers();
    requestHeaders.set('Content-Type', 'application/json')
    requestHeaders.set('Access-Control-Allow-Credentials', 'true');
    requestHeaders.set('Access-Control-Allow-Origin', '*');
    const response = await fetch('http://localhost:8080/games', {
        method: 'POST',
        body: JSON.stringify({ player_names: playerNames }),
        headers: {
            'Access-Control-Allow-Origin': '*',
            'Content-Type': 'application/json',
            "Access-Control-Allow-Methods": "DELETE, POST, GET, OPTIONS"
        }
    });

    if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
    }

    const data: GameSessionResponse = await response.json();
    return data.session_id;
}