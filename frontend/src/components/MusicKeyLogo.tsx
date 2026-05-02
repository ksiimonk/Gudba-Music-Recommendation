type MusicKeyLogoProps = {
  className?: string
  label?: string
}

export function MusicKeyLogo({
  className = '',
  label = 'Gudba Music',
}: MusicKeyLogoProps) {
  return (
    <div className={`music-key-logo ${className}`} aria-label={label}>
      <svg
        aria-hidden="true"
        viewBox="0 0 100 100"
        focusable="false"
      >
        <path d="M39 67.5V31.5L69 24V60.5C66.7 59.3 63.7 58.9 60.6 59.7C54.9 61 51 65 51.8 68.6C52.6 72.1 57.9 73.8 63.6 72.5C68.7 71.3 72.4 67.9 72.5 64.6V16L35.5 25.2V63.5C33.2 62.3 30.2 61.9 27.1 62.7C21.4 64 17.5 68 18.3 71.6C19.1 75.1 24.4 76.8 30.1 75.5C35.1 74.3 38.8 70.9 39 67.5Z" />
      </svg>
    </div>
  )
}
