import React, {useState} from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import HelloWorld from './HelloWorld';
import '@testing-library/jest-dom'
import { act } from 'react-dom/test-utils';

jest.mock('../Env', () => ({
  API_URL: 'http://test-api.com',
}));

describe('HelloWorld', () => {
  it('fetches and displays the message', async () => {
    const mockResponse = { message: 'Hello from the API!' };
    window.fetch = jest.fn().mockResolvedValue({
      json: () => Promise.resolve(mockResponse),
    });

    
    await act(async () => {
        render(<HelloWorld />);
      });

    await waitFor(() => expect(window.fetch).toHaveBeenCalledWith('http://test-api.com/hello_world'));
    expect(screen.getByText(/Hello World/)).toBeInTheDocument();
    expect(screen.getByText(/Hello from the API!/)).toBeInTheDocument();
  });

  it('handles fetch error', async () => {
    const mockError = new Error('fetch error');
    window.fetch = jest.fn().mockRejectedValue(mockError);
    const spy = jest.spyOn(console, 'error').mockImplementation(() => {});

    await act(async () => {
        render(<HelloWorld />);
      });

    await waitFor(() => expect(window.fetch).toHaveBeenCalledWith('http://test-api.com/hello_world'));
    expect(spy).toHaveBeenCalledWith(mockError);

    spy.mockRestore();
  });
});