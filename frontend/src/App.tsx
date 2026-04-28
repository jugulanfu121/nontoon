import './App.css'

function App() {

  const handleInputChange = async (e: React.ChangeEvent<HTMLInputElement, HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0]
      const chunks = createChunks(file, 1024 * 1024)

      console.log("chunks length", chunks.length)

      for (let i = 0; i < chunks.length; i++) {
        const formData = new FormData()

        const f = blobToFile(chunks[i], file.name)

        formData.append("file", f)
        formData.append("uploadId", "def")
        formData.append("chunkIndex", i.toString())
        const response = await fetch("http://localhost:8080/videos/chunks", {
          method: "POST",
          body: formData
        })

        const res = await response.json()
        console.log("result api", res)
      }
    }
  }

  const blobToFile = (theBlob: Blob, fileName: string): File => {
    return new File([theBlob], fileName, {
      type: theBlob.type,
      lastModified: Date.now()
    });
  };

  const createChunks = (file: File, chunkSize: number): Blob[] => {
    const chunks: Blob[] = [];
    for (let start = 0; start < file.size; start += chunkSize) {
      const chunk = file.slice(start, start + chunkSize);
      chunks.push(chunk);
    }
    return chunks;
  };
  return (
    <>
      <input type='file' onChange={(e) => handleInputChange(e)}></input>
    </>
  )
}

export default App
