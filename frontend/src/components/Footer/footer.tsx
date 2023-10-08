// src/components/RecentSong.tsx

import React, { useState, useEffect } from 'react';

interface RecentSongData {
  songName: string;
  songArtist: string;
}

function RecentSong() {
  const [recentSong, setRecentSong] = useState<RecentSongData>({ songName: '', songArtist: '' });

  useEffect(() => {
    const fetchRecentSong = async () => {
      try {
        const response = await fetch('http://localhost:8080/recentSong'); // Replace with your Gin backend URL
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data: RecentSongData = await response.json();
        setRecentSong(data);
      } catch (error) {
        console.error('Error fetching recent song:', error);
      }
    };

    fetchRecentSong();
  }, []);

  return (
    <div>
      Currently listening to: {" "} 
      <span className='highlight'>{recentSong.songName}</span> by {" "}
      <span className='highlight'>{recentSong.songArtist}</span> 
    </div>
  );
}

export default RecentSong;
