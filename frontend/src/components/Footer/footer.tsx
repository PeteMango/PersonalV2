import React, { useState, useEffect } from "react";
import './footer.css'

interface RecentSongData {
  songName: string;
  songArtist: string;
  songLink: string;
}

function RecentSong() {
  const [recentSong, setRecentSong] = useState<RecentSongData | null>(null);

  useEffect(() => {
    const fetchRecentSong = async () => {
      try {
        const response = await fetch("http://localhost:8080/recentSong"); // Replace with your Gin backend URL
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const data: RecentSongData = await response.json();
        setRecentSong(data);
      } catch (error) {
        console.error("Error fetching recent song:", error);
      }
    };

    fetchRecentSong();
  }, []);

  return (
    <div>
      {recentSong ? (
        <div>
          Listening to:{" "}
          <span className="song"><a href={recentSong.songLink} target="_blank">{recentSong.songName}</a></span>
        </div>
      ) : (
        <div>
          <span className="music">Not scrobbling.</span>
        </div>
      )}
    </div>
  );
}

export default RecentSong;
