// import React, { useState, useEffect } from 'react';
// import { useRouter } from 'next/router';

// interface Player {
//   id: number;
//   game_id: string;
//   name: string;
// }

// interface Frame {
//   ID: number;
//   PlayerID: number;
//   FrameNumber: number;
//   Roll1: string;
//   Roll2: string;
//   Roll3: string;
//   Score: number;
// }

// interface BowlingData {
//   players: Player[];
//   frames: Frame[];
// }

// const ScoreBoard: React.FC = () => {
//   const [data, setData] = useState<BowlingData | null>(null);
//   const [loading, setLoading] = useState<boolean>(true);
//   const [error, setError] = useState<string | null>(null);
//   const [scores, setScores] = useState<{ [playerId: number]: { roll1: string; roll2: string; roll3: string } }>({});
//   const router = useRouter();

//   useEffect(() => {
//     if (!router.isReady) return;

//     const { sessionID } = router.query;

//     const fetchData = async () => {
//       try {
//         const response = await fetch(`http://localhost:8080/games/${sessionID}`);
//         if (!response.ok) {
//           throw new Error(`HTTP error! status: ${response.status}`);
//         }
//         const result: BowlingData = await response.json();
//         setData(result);

//         // Initialize scores state with empty values for each player
//         const initialScores: { [playerId: number]: { roll1: string; roll2: string; roll3: string } } = {};
//         result.players.forEach((player) => {
//           initialScores[player.id] = { roll1: '', roll2: '', roll3: '' };
//         });
//         setScores(initialScores);
//       } catch (err: any) {
//         setError(err.message);
//       } finally {
//         setLoading(false);
//       }
//     };

//     fetchData();
//   }, [router.isReady, router.query.sessionID]);

//   const handleInputChange = (playerId: number, roll: 'roll1' | 'roll2' | 'roll3', value: string) => {
//     setScores((prevScores) => ({
//       ...prevScores,
//       [playerId]: {
//         ...prevScores[playerId],
//         [roll]: value,
//       },
//     }));
//   };

//   const handleSubmit = async () => {
//     // Validate that all scores are filled
//     for (const playerId in scores) {
//       const { roll1, roll2, roll3 } = scores[playerId];
//       if (!roll1 && !roll2 && !roll3) {
//         alert('Please fill in all scores before submitting.');
//         return;
//       }
//     }

//     // Prepare the payload
//     const payload = {
//       scores: Object.fromEntries(
//         Object.entries(scores).map(([playerId, rolls]) => [
//           playerId,
//           {
//             frame_number: (data?.frames?.length ?? 0) + 1, // Assuming adding to the next frame
//             roll1: rolls.roll1,
//             roll2: rolls.roll2,
//             roll3: rolls.roll3,
//           },
//         ])
//       ),
//     };

//     try {
//       const response = await fetch(`http://localhost:8080/games/${router.query.sessionID}/frames`, {
//         method: 'POST',
//         headers: { 'Content-Type': 'application/json' },
//         body: JSON.stringify(payload),
//       });

//       if (!response.ok) {
//         throw new Error(`HTTP error! status: ${response.status}`);
//       } else {
//         window.location.reload();
//       }
//     } catch (err) {
//       console.error('Submission error:', err);
//       alert('An error occurred while submitting the scores.');
//     }
//   };

//   if (loading) {
//     return <div>Loading...</div>;
//   }

//   if (error) {
//     return <div>Error: {error}</div>;
//   }

//   return (
//     <div>
//       {data?.players.map((player) => (
//         <div key={player.id}>
//           <h2>{player.name}</h2>
//           <table border={1}>
//             <thead>
//               <tr>
//                 <th>Frame</th>
//                 <th>Roll 1</th>
//                 <th>Roll 2</th>
//                 <th>Roll 3</th>
//                 <th>Score</th>
//               </tr>
//             </thead>
//             <tbody>
//               {data.frames
//                 .filter((frame) => frame.PlayerID === player.id)
//                 .map((frame) => (
//                   <tr key={frame.ID}>
//                     <td>{frame.FrameNumber}</td>
//                     <td>{frame.Roll1}</td>
//                     <td>{frame.Roll2}</td>
//                     <td>{frame.Roll3}</td>
//                     <td>{frame.Score}</td>
//                   </tr>
//                 ))}
//               <tr>
//                 <td>New Frame</td>
//                 <td>
//                   <input
//                     type="text"
//                     value={scores[player.id]?.roll1 || ''}
//                     onChange={(e) => handleInputChange(player.id, 'roll1', e.target.value)}
//                   />
//                 </td>
//                 <td>
//                   <input
//                     type="text"
//                     value={scores[player.id]?.roll2 || ''}
//                     onChange={(e) => handleInputChange(player.id, 'roll2', e.target.value)}
//                   />
//                 </td>
//                 <td>
//                   <input
//                     type="text"
//                     value={scores[player.id]?.roll3 || ''}
//                     onChange={(e) => handleInputChange(player.id, 'roll3', e.target.value)}
//                   />
//                 </td>
//                 <td></td>
//               </tr>
//             </tbody>
//           </table>
//         </div>
//       ))}
//       <button onClick={handleSubmit}>Submit All Scores</button>
//     </div>
//   );
// };

// export default ScoreBoard;


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
  frames: Frame[];
}

const ScoreBoard: React.FC = () => {
  const [data, setData] = useState<BowlingData | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [scores, setScores] = useState<{ [playerId: number]: { roll1: string; roll2: string; roll3: string } }>({});
  const router = useRouter();

  useEffect(() => {
    if (!router.isReady) return;

    const { sessionID } = router.query;

    const fetchData = async () => {
      try {
        const response = await fetch(`http://localhost:8080/games/${sessionID}`);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result: BowlingData = await response.json();
        setData(result);

        // Initialize scores state with empty values for each player
        const initialScores: { [playerId: number]: { roll1: string; roll2: string; roll3: string } } = {};
        result.players.forEach((player) => {
          initialScores[player.id] = { roll1: '', roll2: '', roll3: '' };
        });
        setScores(initialScores);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [router.isReady, router.query.sessionID]);

  const handleInputChange = (playerId: number, roll: 'roll1' | 'roll2' | 'roll3', value: string) => {
    setScores((prevScores) => ({
      ...prevScores,
      [playerId]: {
        ...prevScores[playerId],
        [roll]: value,
      },
    }));
  };

  const handleSubmit = async () => {
    // Validate that all scores are filled
    for (const playerId in scores) {
      const { roll1, roll2, roll3 } = scores[playerId];
      if (!roll1 && !roll2 && !roll3) {
        alert('Please fill in all scores before submitting.');
        return;
      }
    }

    // Prepare the payload
    const payload = {
      scores: Object.fromEntries(
        Object.entries(scores).map(([playerId, rolls]) => [
          playerId,
          {
            frame_number: (data?.frames?.length ?? 0) + 1, // Assuming adding to the next frame
            roll1: rolls.roll1,
            roll2: rolls.roll2,
            roll3: rolls.roll3,
          },
        ])
      ),
    };

    try {
      const response = await fetch(`http://localhost:8080/games/${router.query.sessionID}/frames`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      } else {
        window.location.reload();
      }
    } catch (err) {
      console.error('Submission error:', err);
      alert('An error occurred while submitting the scores.');
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      {data?.players.map((player) => {
        // Calculate total score for the player
        const totalScore = data.frames
          ?.filter((frame) => frame.PlayerID === player.id)
          .reduce((sum, frame) => sum + frame.Score, 0);

        return (
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
                      <td>{frame.FrameNumber}</td>
                      <td>{frame.Roll1}</td>
                      <td>{frame.Roll2}</td>
                      <td>{frame.Roll3}</td>
                      <td>{frame.Score}</td>
                    </tr>
                  ))}
                {/* Total Score Row */}
                <tr>
                  <td colSpan={4}><strong>Total Score</strong></td>
                  <td><strong>{totalScore}</strong></td>
                </tr>
                {/* New Frame Input Row */}
                <tr>
                  <td>New Frame</td>
                  <td>
                    <input
                      type="text"
                      value={scores[player.id]?.roll1 || ''}
                      onChange={(e) => handleInputChange(player.id, 'roll1', e.target.value)}
                    />
                  </td>
                  <td>
                    <input
                      type="text"
                      value={scores[player.id]?.roll2 || ''}
                      onChange={(e) => handleInputChange(player.id, 'roll2', e.target.value)}
                    />
                  </td>
                  <td>
                    <input
                      type="text"
                      value={scores[player.id]?.roll3 || ''}
                      onChange={(e) => handleInputChange(player.id, 'roll3', e.target.value)}
                    />
                  </td>
                  <td></td>
                </tr>
              </tbody>
            </table>
          </div>
        );
      })}
      <button onClick={handleSubmit}>Submit All Scores</button>
    </div>
  );
};

export default ScoreBoard;
