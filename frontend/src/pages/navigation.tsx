import Link from 'next/link';

interface NavigationProps {
  sessionID: string;
}

const Navigation: React.FC<NavigationProps> = ({ sessionID }) => {
  return (
    <nav>
      <Link href="/">Home</Link>
      <Link href="/register">Register</Link>
      <Link href={`/games/${sessionID}`}>ScoreBoard</Link>
    </nav>
  );
};

export default Navigation;