import React, { useState } from 'react';
import {createGameSession} from '../api/game'
import { useRouter } from 'next/router';

const Register: React.FC = () => {
  const [players, setPlayers] = useState<string[]>(Array(5).fill(''));
  const router = useRouter();
  const handleChange = (index: number) => (event: React.ChangeEvent<HTMLInputElement>) => {
    const newPlayers = [...players];
    newPlayers[index] = event.target.value;
    setPlayers(newPlayers);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    
    const filledPlayers = players.filter(name => name.trim() !== '');
    if (filledPlayers.length == 0) {
      alert("number of users has to be at least 1")
      return
    }
    const gameSessionID = await createGameSession(filledPlayers)
   
    router.push(`/games/${gameSessionID}`);
    console.log(gameSessionID)
  };

  return (
    <div style={styles.container}>
      <h2>Register Players</h2>
      <form onSubmit={handleSubmit} style={styles.form}>
        {players.map((player, index) => (
          <div key={index} style={styles.inputGroup}>
            <label htmlFor={`player${index + 1}`} style={styles.label}>
              Player {index + 1}:
            </label>
            <input
              type="text"
              id={`player${index + 1}`}
              value={player}
              onChange={handleChange(index)}
              required = {false}
              style={styles.input}
            />
          </div>
        ))}
        <button type="submit" style={styles.button}>Start Game</button>
      </form>
    </div>
  );
};

const styles = {
  container: {
    fontFamily: 'Arial, sans-serif',
    backgroundColor: '#f2f2f2',
    display: 'flex',
    flexDirection: 'column' as const,
    alignItems: 'center' as const,
    height: '100vh',
    margin: 0,
    paddingTop: '50px',
  },
  form: {
    backgroundColor: '#fff',
    padding: '20px',
    borderRadius: '8px',
    boxShadow: '0 0 10px rgba(0, 0, 0, 0.1)',
    width: '300px',
  },
  inputGroup: {
    marginBottom: '15px',
  },
  label: {
    display: 'block',
    marginBottom: '5px',
    color: 'black',
  },
  input: {
    width: '100%',
    padding: '8px',
    boxSizing: 'border-box' as const,
  },
  button: {
    width: '100%',
    padding: '10px',
    backgroundColor: '#4CAF50',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    cursor: 'pointer',
  },
};

export default Register;


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
//   const [selectedPlayer, setSelectedPlayer] = useState<number | null>(null);
//   const [newScore, setNewScore] = useState({ roll1: '', roll2: '', roll3: '' });
//   const router = useRouter();

//   // useEffect(() => {
//   //   if (!router.isReady) return;

//   //   const { sessionID } = router.query;

//   //   const fetchData = async () => {
//   //     try {
//   //       const response = await fetch(`http://localhost:8080/games/${sessionID}`);
//   //       if (!response.ok) {
//   //         throw new Error(`HTTP error! status: ${response.status}`);
//   //       }
//   //       const result: BowlingData = await response.json();
//   //       setData(result);
//   //     } catch (err: any) {
//   //       setError(err.message);
//   //     } finally {
//   //       setLoading(false);
//   //     }
//   //   };

//   //   fetchData();
//   // }, [router.isReady, router.query.sessionID]);

//   const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
//     const { name, value } = e.target;
//     setNewScore((prev) => ({ ...prev, [name]: value }));
//   };

//   const handleSubmit = async () => {
//     if (selectedPlayer === null) return;

//     const frameNumber =
//       (data?.frames.filter((frame) => frame.PlayerID === selectedPlayer).length || 0) + 1;

//     const newFrame: Frame = {
//       ID: Date.now(), // Temporary ID; replace with actual ID from the backend
//       PlayerID: selectedPlayer,
//       FrameNumber: frameNumber,
//       Roll1: newScore.roll1,
//       Roll2: newScore.roll2,
//       Roll3: newScore.roll3,
//       Score: parseInt(newScore.roll1) + parseInt(newScore.roll2) + parseInt(newScore.roll3),
//     };

//     // Here you can send newFrame to your API
//     console.log('Submitting new frame:', newFrame);

//     // Reset input fields
//     setNewScore({ roll1: '', roll2: '', roll3: '' });
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
//             </tbody>
//           </table>
//         </div>
//       ))}

//       <div>
//         <h2>Add New Score</h2>
//         <select
//           onChange={(e) => setSelectedPlayer(parseInt(e.target.value))}
//           value={selectedPlayer || ''}
//         >
//           <option value="" disabled>
//             Select Player
//           </option>
//           {data?.players.map((player) => (
//             <option key={player.id} value={player.id}>
//               {player.name}
//             </option>
//           ))}
//         </select>

//         {selectedPlayer !== null && (
//           <div>
//             <input
//               type="text"
//               name="roll1"
//               placeholder="Roll 1"
//               value={newScore.roll1}
//               onChange={handleInputChange}
//             />
//             <input
//               type="text"
//               name="roll2"
//               placeholder="Roll 2"
//               value={newScore.roll2}
//               onChange={handleInputChange}
//             />
//             <input
//               type="text"
//               name="roll3"
//               placeholder="Roll 3"
//               value={newScore.roll3}
//               onChange={handleInputChange}
//             />
//             <button onClick={handleSubmit}>Submit Score</button>
//           </div>
//         )}
//       </div>
//     </div>
//   );
// };

// export default ScoreBoard;
