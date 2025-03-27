interface GameSessionResponse {
    session_id: string;
  }
  
export async function createGameSession(playerNames: string[]): Promise<string> {
    const response = await fetch('http://localhost:8080/games', {
        method: 'POST',
        body: JSON.stringify({ player_names: playerNames }),
        headers: {
            'Content-Type': 'application/json',
        }
    });

    if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
    }

    const data: GameSessionResponse = await response.json();
    return data.session_id;
}

export interface Player {
    id: number;
    game_id: string;
    name: string;
  }
  
  export  interface Frame {
    ID: number;
    PlayerID: number;
    FrameNumber: number;
    Roll1: string;
    Roll2: string;
    Roll3: string;
    Score: number;
  }
  
export  interface BowlingData {
    players: Player[];
    frames: Frame[];
  }

export async function getGameInfo(sessionID: string): Promise<BowlingData> {
    const response = await fetch(`http://localhost:8080/games/${sessionID}`);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const result: BowlingData = await response.json();
    return result
}