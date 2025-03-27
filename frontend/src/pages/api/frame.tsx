
/* eslint-disable  @typescript-eslint/no-explicit-any */
export async function createFrames(sessionID: string, frameData: any) {
    const response = await fetch(`http://localhost:8080/games/${sessionID}/frames`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(frameData),
      })

      if (!response.ok) {
        throw new Error(`Error: ${response.status} ${response.statusText}`);
    }
}