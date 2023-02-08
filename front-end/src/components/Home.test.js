import React from 'react';
import { render } from '@testing-library/react';
import Home from './Home';
import '@testing-library/jest-dom'

describe('Home', () => {
    it('renders logo', () => {
        const { getByAltText } = render(<Home />);
        const logoElement = getByAltText('logo');
        expect(logoElement).toBeInTheDocument();
    });

    it('renders "Bonjour Monde !" text', () => {
        const { getByText } = render(<Home />);
        const textElement = getByText(/Bonjour Monde !/i);
        expect(textElement).toBeInTheDocument();
    });

    it('renders "Learn React" link', () => {
        const { getByText } = render(<Home />);
        const linkElement = getByText(/Learn React/i);
        expect(linkElement).toBeInTheDocument();
    });

    it('renders "App-header" class', () => {
        const { container } = render(<Home />);
        const headerElement = container.querySelector('.App-header');
        expect(headerElement).toBeInTheDocument();
    });
});