import { useState } from 'react'
import './App.css'

function App() {

  const [totalChunks, setTotalChunks] = useState(0)
  const [totalChunksUploaded, setTotalChunksUploaded] = useState(0)
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
        setTotalChunksUploaded(i)
        console.log("result api", res)
      }

      await completeUpload("def", file.name, chunks.length)
    }
  }

  const completeUpload = async (uploadId: string, fileName: string, totalChunks: number) => {
    const formData = new FormData()
    formData.append("uploadId", uploadId)
    formData.append("filename", fileName)
    formData.append("totalChunks", totalChunks.toString())

    await fetch("http://localhost:8080/videos/merge", {
      method: "POST",
      body: formData
    })
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
    setTotalChunks(chunks.length)
    return chunks;
  };
  return (
    <>
      <input type='file' onChange={(e) => handleInputChange(e)}></input>

      <br />

      <p>Total chunks uploaded: {totalChunksUploaded} out of {totalChunks}</p>
    </>
  )
}

export default App
