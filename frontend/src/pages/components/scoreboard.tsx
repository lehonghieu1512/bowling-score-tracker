import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';

interface Player {
  id: number;
  game_id: string;
  name: string;
}

interface Frame {
  ID: number;
  PlayerID: number;
  FrameNumber: number;
  Roll1: string;
  Roll2: string;
  Roll3: string;
  Score: number;
}

interface BowlingData {
  players: Player[];
  frames: Frame[] | null;
}

const ScoreBoard: React.FC = () => {
  const [data, setData] = useState<BowlingData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  useEffect(() => {
    const fetchData = async () => {
      try {
        const { sessionID } = router.query;
        const response = await fetch(`http://localhost:8080/games/${sessionID}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      {data?.players.map((player) => (
        <div key={player.id}>
          <h2>{player.name}</h2>
          <table border={1}>
            <thead>
              <tr>
                <th>Frame</th>
                <th>Roll 1</th>
                <th>Roll 2</th>
                <th>Roll 3</th>
                <th>Score</th>
              </tr>
            </thead>
            <tbody>
              {data.frames
                ?.filter((frame) => frame.PlayerID === player.id)
                .map((frame) => (
                  <tr key={frame.ID}>
                    <td>{frame.FrameNumber + 1}</td>
                    <td>{frame.Roll1}</td>
                    <td>{frame.Roll2}</td>
                    <td>{frame.Roll3}</td>
                    <td>{frame.Score}</td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      ))}
    </div>
  );
};

export default ScoreBoard;
