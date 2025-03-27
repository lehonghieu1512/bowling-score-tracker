import '@/styles/globals.css';
import type { AppProps } from 'next/app';
// import Navigation from '../pages/navigation'; // Adjust the path as necessary

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      {/* <Navigation /> */}
      <Component {...pageProps} />
    </>
  );
}

export default MyApp;