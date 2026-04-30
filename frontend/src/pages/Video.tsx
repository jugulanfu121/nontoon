import { MediaControlBar, MediaController, MediaFullscreenButton, MediaMuteButton, MediaPlaybackRateButton, MediaPlayButton, MediaSeekBackwardButton, MediaSeekForwardButton, MediaTimeDisplay, MediaTimeRange, MediaVolumeRange } from 'media-chrome/react'
import { MediaRenditionMenu, MediaRenditionMenuButton } from 'media-chrome/react/menu'
import ReactPlayer from 'react-player'

const Video = () => {
    return (
        <MediaController
            style={{
                width: "100%",
                aspectRatio: "16/9",
            }}
        >
            <ReactPlayer
                slot="media"
                src="http://localhost:8080/assets/def/master.m3u8"
                controls={false}
                style={{
                    width: "100%",
                    height: "100%",
                }}
            ></ReactPlayer>
            <MediaRenditionMenu anchor="auto" />
            <MediaControlBar>
                <MediaPlayButton />
                <MediaSeekBackwardButton seekOffset={10} />
                <MediaSeekForwardButton seekOffset={10} />
                <MediaTimeRange />
                <MediaTimeDisplay showDuration />
                <MediaRenditionMenuButton />
                <MediaMuteButton />
                <MediaVolumeRange />
                <MediaPlaybackRateButton />
                <MediaFullscreenButton />
            </MediaControlBar>
        </MediaController>
    )
}

export default Video