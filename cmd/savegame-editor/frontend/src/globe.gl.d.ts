declare module 'globe.gl' {
  import { Object3D, Material } from 'three';

  interface GlobeInstance {
    (element: HTMLElement): GlobeInstance;

    // Globe appearance
    globeImageUrl(url: string): GlobeInstance;
    globeMaterial(material: Material): GlobeInstance;
    backgroundColor(color: string): GlobeInstance;
    showAtmosphere(show: boolean): GlobeInstance;
    showGraticules(show: boolean): GlobeInstance;

    // Points layer
    pointsData(data: object[]): GlobeInstance;
    pointLat(accessor: string | ((d: any) => number)): GlobeInstance;
    pointLng(accessor: string | ((d: any) => number)): GlobeInstance;
    pointColor(accessor: string | ((d: any) => string)): GlobeInstance;
    pointRadius(accessor: string | ((d: any) => number)): GlobeInstance;
    pointAltitude(accessor: string | ((d: any) => number)): GlobeInstance;
    pointLabel(accessor: string | ((d: any) => string)): GlobeInstance;

    // Camera
    pointOfView(pov: { lat?: number; lng?: number; altitude?: number }, transitionMs?: number): GlobeInstance;

    // Size
    width(width: number): GlobeInstance;
    height(height: number): GlobeInstance;

    // Scene access
    scene(): Object3D;
    renderer(): { dispose(): void };
    controls(): { dispose(): void };

    // Lifecycle
    pauseAnimation(): GlobeInstance;
    resumeAnimation(): GlobeInstance;

    _destructor?(): void;
  }

  export default function Globe(): GlobeInstance;
}
