import React, {useCallback, useEffect, useRef, useState} from 'react';

function ImageDownloader() {
    const [urls, setUrls] = useState([]);
    const [progress, setProgress] = useState({});
    const [urlInput, setUrlInput] = useState('');
    const [images, setImages] = useState([]);
    const [page, setPage] = useState(1);
    const loadingRef = useRef(null);
    const isLoading = useRef(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [currentImage, setCurrentImage] = useState(null);

    const openModal = (image) => {
        setCurrentImage(image);
        setIsModalOpen(true);
    };

    const closeModal = () => {
        setIsModalOpen(false);
    };


    function getExtension(filename) {
        const parts = filename.split('.');
        return parts[parts.length - 1];
    }

    const loadImages = useCallback(async (pageNum) => {
        if (!isLoading.current) {
            isLoading.current = true;
            const response = await fetch(`http://localhost:8080/api/v1/images?page=${pageNum}&pageSize=5`);
            const data = await response.json();
            if (data.data) {
                setImages((prevImages) => {
                    const imagesMap = new Map(prevImages.map(img => [img.id, img]));
                    data.data.forEach(img => {
                        const newImage = {
                            id: img.id,
                            url: `http://localhost:8080/public/${img.id}_mini.` + getExtension(img.name),
                            urlReal: `http://localhost:8080/public/${img.id}.` + getExtension(img.name),
                            created_at: img.created_at,
                        };
                        imagesMap.set(img.id, newImage);
                    });

                    return Array.from(imagesMap.values())
                        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
                });
                setPage(pageNum + 1);
                isLoading.current = false;
            }
        }
    }, []);

    const deleteImage = async (imageId) => {
        try {
            const response = await fetch(`http://localhost:8080/api/v1/images/${imageId}`, {
                method: 'DELETE',
            });
            if (!response.ok) {
                throw new Error('Error al eliminar la imagen');
            }
            setImages(images => images.filter(image => image.id !== imageId));
        } catch (error) {
            console.error(error);
        }
    };

    useEffect(() => {
        loadImages(page).catch(console.error);
    }, [loadImages, page]);

    useEffect(() => {
        const observer = new IntersectionObserver((entries) => {
            if (entries[0].isIntersecting) {
                loadImages(page).catch(console.error);
            }
        }, {
            rootMargin: '100px',
        });

        if (loadingRef.current) {
            observer.observe(loadingRef.current);
        }

        return () => observer.disconnect();
    }, [page, loadImages]);

    const generateUUID = () => {
        let dt = new Date().getTime();
        return 'xxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
            const r = (dt + Math.random() * 16) % 16 | 0;
            dt = Math.floor(dt / 16);
            return (c === 'x' ? r : ((r & 0x3) | 0x8)).toString(16);
        });
    };


    const addUrl = () => {
        if (urlInput && urls.filter(url => url.url === urlInput).length === 0) {
            const uuid = generateUUID();
            setUrls([...urls, {url: urlInput, id: uuid}]);
            setProgress({...progress, [uuid]: {progress: 0, completed: false}});
            setUrlInput('');
        }
    };

    const startDownload = () => {
        const ws = new WebSocket('ws://localhost:8080/ws');
        ws.onopen = () => {
            ws.send(JSON.stringify({command: "start_download", data: urls}));
        };

        ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.event === 'progress') {
                setProgress(prevProgress => {
                    return {
                        ...prevProgress, [data.id]: {
                            progress: data.progress, completed: data.completed, error: data.error || ''
                        }
                    };
                });
            } else if (data.event === 'final') {
                setTimeout(function(){
                    isLoading.current = false;
                    loadImages(1).catch(console.error)
                }, 2000);
            }
        };

        ws.onclose = (event) => {
            console.log('Conexión cerrada', event);
        };

        ws.onerror = (error) => {
            console.log('WebSocket Error', error);
        };
    };

    const removeUrl = (id) => {
        setUrls(urls.filter(url => url.id !== id));
        const newProgress = {...progress};
        delete newProgress[id];
        setProgress(newProgress);
    };

    return (<div className="flex flex-col justify-center items-center min-h-screen p-4">
            <h1 className="text-3xl font-bold mb-6">Descarga de Imágenes</h1>
            <div className="mb-8 w-full w-max">
                <input
                    type="text"
                    value={urlInput}
                    onChange={(e) => setUrlInput(e.target.value)}
                    placeholder="URL de la imagen"
                    className="w-full border-2 border-gray-300 rounded py-2 px-4 mb-4"
                />
                <div className="flex justify-center gap-4">
                    <button onClick={addUrl}
                            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded transition duration-300 ease-in-out">Añadir
                        URL
                    </button>
                    <button onClick={startDownload}
                            className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-6 rounded transition duration-300 ease-in-out">Iniciar
                        Descarga
                    </button>
                </div>
            </div>

            <div className="w-full px-4 mb-8">
                {urls.map((url) => (<div key={url.id} className="mb-4 flex items-center justify-between">
                        <div className="flex-1">
                            <p className="text-gray-700 truncate">{url.url}</p>
                            {!!progress[url.id]?.error && (<p className="text-red-500">{progress[url.id]?.error}</p>)}
                            <div className="w-full bg-gray-300 rounded-full h-6">
                                <div
                                    className={`h-6 bg-blue-500 rounded-full transition-all duration-500 ease-in-out ${progress[url.id]?.completed ? 'bg-green-500' : ''}`}
                                    style={{width: `${progress[url.id]?.progress.toFixed(2) || 0}%`}}
                                >
                                    <span
                                        className="text-sm text-white pl-2">{`${progress[url.id]?.progress.toFixed(2) || 0}%`}</span>
                                </div>
                            </div>
                        </div>
                        <button
                            onClick={() => removeUrl(url.id)}
                            className="ml-4 bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 rounded transition duration-300 ease-in-out"
                            disabled={!progress[url.id]?.completed}
                        >
                            X
                        </button>
                    </div>))}
            </div>

            <h2 className="text-2xl font-bold mb-4">Galería de Imágenes</h2>
            <div className="flex flex-wrap justify-center gap-4 pb-4">
                {images.map((image) => (<div key={image.id} className="w-48 relative group">
                        <div className="aspect-w-1 aspect-h-1">
                            <img src={image.url} alt="" className="rounded-lg shadow-md object-cover w-full h-full" onClick={() => openModal(image)}/>
                            <button
                                onClick={() => deleteImage(image.id)}
                                className="absolute top-2 right-2 bg-black bg-opacity-50 hover:bg-opacity-70 text-white font-bold py-2 px-2 rounded-full transition duration-300 ease-in-out opacity-0 group-hover:opacity-100"
                                style={{
                                    width: '30px',
                                    height: '30px',
                                    display: 'flex',
                                    alignItems: 'center',
                                    justifyContent: 'center'
                                }}
                            >
                                X
                            </button>
                        </div>
                    </div>))}
            </div>
            {isModalOpen && (<div
                    className="fixed top-0 left-0 w-full h-full bg-black bg-opacity-70 flex justify-center items-center z-50"
                    onClick={closeModal}>
                    <img src={currentImage?.urlReal} alt="" className="max-w-full max-h-full z-50"/>
                </div>)}
            <div ref={loadingRef} className="h-10"></div>
        </div>);

}

export default ImageDownloader;
